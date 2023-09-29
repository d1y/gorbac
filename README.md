# gorbac

使用到的库: https://github.com/storyicon/grbac

参考: https://github.com/saltbo/zpan/blob/master/internal/pkg/middleware/auth_rbac.yml

```yaml
# 默认任何路由都不需要角色
- id: 0
  host: "*"
  path: "**"
  method: "*"
  allow_anyone: true

# 超级管理员
- id: 42
  host: "*"
  path: "/api/user/admin/**"
  method: "*"
  authorized_roles:
    - "superadmin"

# 审核员和超级管理员
- id: 12
  host: "*"
  path: "/api/**"
  method: "*"
  authorized_roles:
    - "review"
    - "superadmin"
    - "admin"
```