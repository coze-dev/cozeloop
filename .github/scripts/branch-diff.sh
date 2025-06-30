#!/bin/bash

# 解析命令行参数
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --target-branch=*)
      target_branch="${1#*=}"
      ;;
    --changed-files-path=*)
      changed_files_path="${1#*=}"
      ;;
    *)
      echo "Unknown parameter: $1"
      exit 1
  esac
  shift
done

# 如果未提供目标分支,使用 main 分支
if [ -z "$target_branch" ]; then
  target_branch="main"
fi

# 获取当前分支名称
current_branch=$(git rev-parse --abbrev-ref HEAD)

# # 切换到目标分支
# git checkout $target_branch > /dev/null 2>&1

# # 统计落后提交数
# behind_commits=$(git rev-list --left-right --count $current_branch...$target_branch | awk '{print $2}')

# # 切换回当前分支
# git checkout $current_branch > /dev/null 2>&1

# 获取变更文件名称
changed_files_array=($(git diff --name-only $target_branch -- 'common/*' 'frontend/*'))
changed_files_formatted=$(printf '"%s",' "${changed_files_array[@]}")
changed_files_formatted="[${changed_files_formatted%,}]"

# 将输出写入文件
if [ -n "$changed_files_path" ]; then
  echo "$changed_files_formatted" > "$changed_files_path"
else
  echo "changed_files=$changed_files_formatted"
fi