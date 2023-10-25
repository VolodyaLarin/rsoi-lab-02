{{- define "service.template" }}
apiVersion: v1
kind: Service
metadata:
  name: {{.ctx.Release.Name}}-{{.service.name}}-srv
spec:
  selector:
    app: {{.ctx.Release.Name}}-{{.service.name}}
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: NodePort
{{- end }}


  