#!/bin/bash

# 如果未提供目标分支,使用 main 分支
if [ -z "$1" ]; then
  target_branch="origin/main"
else
  target_branch=$1
fi

diff_list=$(git diff --name-only "$target_branch" -- 'common/*' 'frontend/*')
convert_to_json_array() {
  local input_string="$1"
  local output_array="["

  # 使用 IFS 变量临时更改字段分隔符为换行符
  old_ifs=$IFS
  IFS=$'\n'

  # 将字符串分割为单词并添加到数组中
  for word in $input_string; do
    output_array+="\"$word\","
  done

  # 移除最后一个逗号并关闭数组
  output_array=${output_array%,}"]"

  # 恢复原始 IFS 值
  IFS=$old_ifs

  echo "$output_array"
}

echo "$(convert_to_json_array "$diff_list")"