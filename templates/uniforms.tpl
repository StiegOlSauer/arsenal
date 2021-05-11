private _uniforms_hash = createHashMap;

{{- range $u_camo, $u_array := . }}
private _{{$u_camo}} = [
    {{- $printed_once := false}}
    {{- range $i, $u := $u_array }}
    {{- if $printed_once }}, {{end}}"{{$u}}"{{ $printed_once = true}}
    {{- end -}}
];

_uniforms_hash set ["{{ $u_camo }}", _{{ $u_camo }}];
{{ end }}

_uniforms_hash