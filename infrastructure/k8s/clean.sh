yq -i '(select(.data[]) | .data[]) style="double"' ./dist/0002-idlegame.k8s.yaml
