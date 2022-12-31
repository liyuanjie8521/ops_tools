# Git常用命令使用

## 常用命令汇总

### 查看凭据相关命令

1. 查看凭据
```shell
git config --list --show-origin
git config --show-origin --get credential.helper

```

## FAQ汇总
1. 报错:`git push the requested url returned error 403`

   - 解决办法   
   - https操作需要使用用户名密码进行访问

- 清除使用密码缓存配置
```shell
git config --local --unset credential.helper
git config --global --unset credential.helper
git config --system --unset credential.helper
```

- 配置为缓存密码
```shell
git config --global credential.helper store
```