{{/*
Expand the name of the chart.
*/}}
{{- define "updatemonitor.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "updatemonitor.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "updatemonitor.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "updatemonitor.labels" -}}
helm.sh/chart: {{ include "updatemonitor.chart" . }}
{{ include "updatemonitor.selectorLabels.front" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "updatemonitor.selectorLabels.agent" -}}
app.kubernetes.io/name: {{ include "updatemonitor.name" . }}-agent
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
{{- define "updatemonitor.selectorLabels.server" -}}
app.kubernetes.io/name: {{ include "updatemonitor.name" . }}-server
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
{{- define "updatemonitor.selectorLabels.front" -}}
app.kubernetes.io/name: {{ include "updatemonitor.name" . }}-front
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "updatemonitor.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "updatemonitor.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
