# koop

File based Kubernetes Operation Tool

## Usage

**Connect Kubernetes Cluster**

Put kubeconfig files at `$HOME/.koops/cluster-[CLUSTER-NAME].yaml`

**Pull Resources**

```shell
koop pull [CLUSTER-NAME] [NAMESPACE] [KIND] [NAME]
```

Use `-` for wildcard matching

**Push Resource**

```shell
koops push [CLUSTER-NAME] [NAMESPACE] [KIND] [NAME]
```

Use `-` for wildcard matching

## Credits

Guo Y.K., MIT License