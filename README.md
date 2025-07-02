# go-pideaapi
Go package for work with PrivacyIdea API

Currently added functions:
* PrivacyIdea API request: get Pidea API Token for given user(POST) - returns string(API token)
    * `https://{{base_url}}/auth?username={{username}}&password={{password}}`
* PrivacyIdea API request: validate check using realm, user(POST) - returns bool(Check OTP only, without PIN)
    * `https://{{base_url}}/validate/check?user=***&realm=***&otponly=1&pass=***&serial=******`
* PrivacyIdea API request: get users't token serial using using user(GET) - returns string(Token's serial)
    * `https://{{base_url}}/token?realm=***&user=***`
* PrivacyIdea API request: get list of Pidea users of given realm(GET) - list of users' info
    * `https://{{base_url}}/user?realm=***`