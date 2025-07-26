// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package pswd

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
)

// Argon2id 参数
type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// 默认的 Argon2id 参数
var defaultArgon2Params = &argon2Params{
	memory:      64 * 1024, // 64MB
	iterations:  3,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

// HashPassword 使用 Argon2id 算法对密码进行哈希
func HashPassword(password string) (string, error) {
	p := defaultArgon2Params

	// 生成随机盐值
	salt := make([]byte, p.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// 使用 Argon2id 算法计算哈希值
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// 编码为 base64 格式
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式：$argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解析编码后的哈希字符串
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errorx.New("invalid hash format")
	}

	var p argon2Params
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	p.saltLength = uint32(len(salt))

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	p.keyLength = uint32(len(decodedHash))

	// 使用相同的参数和盐值计算哈希值
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// 比较计算得到的哈希值与存储的哈希值
	return subtle.ConstantTimeCompare(decodedHash, computedHash) == 1, nil
}
