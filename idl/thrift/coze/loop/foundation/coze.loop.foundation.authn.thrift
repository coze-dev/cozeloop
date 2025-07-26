namespace go coze.loop.foundation.authn

include "../../../base.thrift"
include "domain/authn.thrift"

struct CreatePersonalAccessTokenRequest {
    1: required string name // PAT名称
    2: optional i64 expire_at // PAT自定义过期时间unix，秒
    3: optional DurationDay duration_day // PAT用户枚举过期时间 1、30、60、90、180、365、permanent

    255: optional base.Base Base
}

typedef string DurationDay(ts.enum="true")
const DurationDay DurationDay_Day1 = "1"
const DurationDay DurationDay_Day30 = "30"
const DurationDay DurationDay_Day60 = "60"
const DurationDay DurationDay_Day90 = "90"
const DurationDay DurationDay_Day180 = "180"
const DurationDay DurationDay_Day365 = "365"
const DurationDay DurationDay_Permanent = "permanent"

struct CreatePersonalAccessTokenResponse {
    1: optional authn.PersonalAccessToken personal_access_token
    2: optional string token    // PAT token 明文

    255: optional base.BaseResp BaseResp
}

struct DeletePersonalAccessTokenRequest {
    1: required i64 id (api.path="id", api.js_conv="true")// PAT id

    255: optional base.Base Base
}

struct DeletePersonalAccessTokenResponse {

    255: optional base.BaseResp BaseResp
}

struct GetPersonalAccessTokenRequest {
    1: required i64 id (api.path="id", api.js_conv="true") // PAT Id

    255: optional base.Base Base
}

struct GetPersonalAccessTokenResponse {
    1: optional authn.PersonalAccessToken personal_access_token

    255: optional base.BaseResp BaseResp
}

struct ListPersonalAccessTokenRequest {
    1: optional i32 page_size (api.query='page_size', vt.not_nil='true', vt.gt='0', vt.le='100') // per page size
    2: optional i32 page_number (api.query='page_number', vt.not_nil='true', vt.gt='0')          // page number

    255: optional base.Base Base
}
struct ListPersonalAccessTokenResponse {
    1: optional list<authn.PersonalAccessToken> personal_access_tokens

    255: optional base.BaseResp BaseResp
}

struct UpdatePersonalAccessTokenRequest {
    1: required i64 id (api.path = "id", api.js_conv="true") // PAT Id
    2: string name // PAT 名称

    255: optional base.Base Base
}

struct UpdatePersonalAccessTokenResponse {
    255: optional base.BaseResp BaseResp
}

struct VerifyTokenRequest {
    1: required string token

    255: optional base.Base Base
}
struct VerifyTokenResponse {
    1: optional bool valid
    2: optional string user_id

    255: optional base.BaseResp BaseResp
}

service AuthNService {
    // OpenAPI PAT管理
    CreatePersonalAccessTokenResponse CreatePersonalAccessToken(1: CreatePersonalAccessTokenRequest req) (api.post='/api/auth/v1/personal_access_tokens')
    DeletePersonalAccessTokenResponse DeletePersonalAccessToken(1: DeletePersonalAccessTokenRequest req) (api.delete='/api/auth/v1/personal_access_tokens/:id')
    UpdatePersonalAccessTokenResponse UpdatePersonalAccessToken(1: UpdatePersonalAccessTokenRequest req) (api.put='/api/auth/v1/personal_access_tokens/:id')
    GetPersonalAccessTokenResponse GetPersonalAccessToken(1: GetPersonalAccessTokenRequest req) (api.get='/api/auth/v1/personal_access_tokens/:id')
    ListPersonalAccessTokenResponse ListPersonalAccessToken(1: ListPersonalAccessTokenRequest req) (api.post='/api/auth/v1/personal_access_tokens/list')

    // 验证token是否有效
    VerifyTokenResponse VerifyToken(1: VerifyTokenRequest req)
}
