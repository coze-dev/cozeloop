#!/usr/bin/env bash
set -e pipefail

SCRIPT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_PATH/../../.."
PROJECT_PATH=$(pwd)
THRIFT_PATH=${PROJECT_PATH}/idl/thrift
gen_module="github.com/coze-dev/coze-loop/backend"

# Define output directories
OUTPUT_DIR=${PROJECT_PATH}/output
KITEX_GEN_DIR=${PROJECT_PATH}/backend/kitex_gen
LOOP_GEN_DIR=${PROJECT_PATH}/backend/loop_gen

# step 1: generate go code from thrift, workdir=output
all_thrifts=($(find $THRIFT_PATH -name "*.thrift" -type f))
mkdir -p ${OUTPUT_DIR} && cd ${OUTPUT_DIR}
for thrift in ${all_thrifts[@]}; do
  rel_path=${thrift#$THRIFT_PATH/}
  dir_path=$(dirname "$rel_path")
  dot_path=${dir_path//\//.}
  thrift_name=$(basename "$thrift" .thrift)
  
  if [ "$thrift_name" = "$dot_path" ]; then
    echo "Processing matching thrift file: $thrift"
    if grep -q "^service" $thrift; then
      echo "thrift service file: $thrift"
      cloudwego_kitex -streamx -thrift ignore_initialisms=false \
        -module ${gen_module} \
        -thrift-plugin validator \
        -thrift=nil_safe \
        -deep-copy-api=true \
        ${thrift}
      cloudwego_hz model --mod=${gen_module} --idl=${thrift} --model_dir=kitex_gen -t=nil_safe -t=streamx -t=thrift_streaming -t=ignore_initialisms=false -t=gen_setter -t=gen_deep_equal -t=compatible_names -t=frugal_tag
    fi
  fi
done

# step 2: loopgen
# copy localstream.go
mkdir -p ${PROJECT_PATH}/output/loop_gen/infra/kitex/localstream
cp -r ${SCRIPT_PATH}/tmpl/*.go ${PROJECT_PATH}/output/loop_gen/infra/kitex/localstream/
loopgen --gomod ${gen_module} \
  --idl-dir ${THRIFT_PATH} \
  --package-prefix "lo" \
  --import-prefix "loop_gen" \
  --local-stream-import-path "github.com/coze-dev/coze-loop/backend/loop_gen/infra/kitex/localstream" \
  -o ${OUTPUT_DIR}

# step 3: move generated files to backend/kitex_gen
rm -rf ${KITEX_GEN_DIR}
rm -rf ${LOOP_GEN_DIR}
mv ${OUTPUT_DIR}/{kitex_gen,loop_gen} ${PROJECT_PATH}/backend/
rm -rf ${OUTPUT_DIR}

echo "Checking for changes in backend/kitex_gen and backend/loop_gen..."
# Check for both tracked and untracked changes
if git diff --quiet ${KITEX_GEN_DIR}/ ${LOOP_GEN_DIR}/ 2>/dev/null &&
  [ -z "$(git ls-files --others --exclude-standard ${KITEX_GEN_DIR}/ ${LOOP_GEN_DIR}/)" ]; then
  echo "No changes detected in kitex_gen and loop_gen, skipping commit."
  exit 0
fi

cd ${PROJECT_PATH}

PUSH_REMOTE=true
if [ "${NO_PUSH_REMOTE}" = true ]; then
    PUSH_REMOTE=false
fi

if [ "$PUSH_REMOTE" = false ]; then
  echo "Skipping git operations as --push_remote was not specified"
  exit 0
fi

echo "Committing backend kitex_gen and loop_gen..."
git add --all
git config user.name "${ACTOR}"
git config user.email "${ACTOR}@bytedance.com"
if [[ -n $(git status --porcelain) ]]; then
  COMMIT=$(git log -1 --pretty=format:"%h")
  git commit -F- <<EOF
ci_trigger: ${ACTOR} ${HEAD_MESSAGE}

author: @${ACTOR}
generate: gen kitex code for ${COMMIT}
EOF
  git push -f origin HEAD:"${HEAD_BRANCH}"
else
  echo "No changes to commit, skipping commit step."
fi
