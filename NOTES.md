

```
helm upgrade --install kindex oci://quay.io/kubosa/charts/kindex:0.1.0-snapshot \
  --namespace kindex --create-namespace \
  --set ingress.host=kindex.ingress.kubo4.mbp \
  --set ingress.class=nginx \
  --set server.clusterName=kubo4

```

```
helm -n kindex uninstall kindex

```



```
helm upgrade --install kindex oci://quay.io/kubotal/charts/kindex:0.1.0-snapshot \
  --namespace kindex --create-namespace \
  --set ingress.host=kindex.ingress.kubo4.mbp \
  --set ingress.class=nginx \
  --set server.clusterName=kubo4

```
