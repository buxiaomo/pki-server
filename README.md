# PKI Server

This project is a management K8S cluster certificate

## Add Project

```shell
curl --location -g --request POST 'http://pki.example.com/v1/pki/project' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "kube-ansible",
    "env": "dev",
    "year": 1
}'
```

## Remove Project

```shell
curl --location -g --request DELETE 'http://pki.example.com/v1/pki/project/kube-ansible/dev'
```

## Renewal Project

```shell
curl --location -g --request PUT 'http://pki.example.com/v1/pki/project/kube-ansible/dev/1'
```

## Renewal Project

```shell
curl --location -g --request PUT 'http://pki.example.com/v1/pki/project/kube-ansible/dev/1'
```

## Get All Project

```shell
curl --location -g --request GET 'http://pki.example.com/v1/pki/project'
```

## Get Project Certificate file

```shell
curl --location -g --request GET 'http://pki.example.com/v1/pki/project/kube-ansible/dev/ca.crt'
curl --location -g --request GET 'http://pki.example.com/v1/pki/project/kube-ansible/dev/ca.crt.sha256sum'
```