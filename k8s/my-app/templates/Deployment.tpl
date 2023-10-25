{{- define "deployment.template" }}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .ctx.Release.Name }}-{{.service.name}}-dep
  labels:
    app: {{ .ctx.Release.Name }}-{{.service.name}}
spec:
  replicas: {{.service.replicaCount}}
  selector:
    matchLabels:
      app: {{ .ctx.Release.Name }}-{{.service.name}}
  template:
    metadata:
      name: {{ .ctx.Release.Name }}-{{.service.name}}
      labels:
        app: {{ .ctx.Release.Name }}-{{.service.name}}
    spec:
      containers:
        - name: {{ .ctx.Release.Name }}-{{.service.name}}
          image: {{.service.container}}
          imagePullPolicy: IfNotPresent
          env:
            {{- range $k, $v := .service.env}}
            - name: {{$k | quote}}
            {{- $fixedadr := $v}}
            {{- range $serviceName,$v := $.services}}
            {{- $fixedadr = ($fixedadr | replace (printf "://%s/" $serviceName) (printf "://%s-%s-srv/" $.ctx.Release.Name $serviceName) )}}
            {{- $fixedadr = ($fixedadr | replace (printf "host=%s" $serviceName) (printf "host=%s-%s-srv" $.ctx.Release.Name $serviceName) )}}
            {{- end}}
              value: {{$fixedadr | quote}}
            {{- end }}
      restartPolicy: Always
{{- end}}