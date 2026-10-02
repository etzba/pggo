# pggo
Simple http service written in go


### install

To set it up in kubernetes:
```sh
helm install pggo chart/ -n pggo --create-namespace
helm upgrade --install pggo chart/ -n pggo 
```
