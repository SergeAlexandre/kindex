

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
helm upgrade --install kindex oci://quay.io/kubosa/charts/kindex:0.1.0-snapshot \
  --namespace kindex --create-namespace \
  --set ingress.host=kindex.ingress.kubo1.mbp \
  --set ingress.class=nginx \
  --set ingress.passthrough=true \
  --set server.clusterName=kubo1 \
  --set server.tls=true \
  --set server.certificateIssuer=cluster-odp \
  --set image.pullPolicy=Always

```
