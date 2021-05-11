private _arsenal_optics = createHashMap;


{{- range $o_type, $o_typesToMountsMap := . }}
private _optics_{{ $o_type }} = createHashMap;
{{- range $o_mount, $o_array := . }}

// Block: optics type: {{ $o_type }}, mount type: {{$o_mount}}
private _{{$o_type}}_{{$o_mount}}_nocamo = [
    {{- $printed_once := false}}
    {{- range $i, $o := $o_array }}
    {{- if eq $o.Camoflage "nocamo" }}{{ if $printed_once }}, {{end}}"{{$o.Name}}"{{ $printed_once = true}}{{ end }}
    {{- end -}}
];
private _{{$o_type}}_{{$o_mount}}_black = [
    {{- $printed_once := false}}
    {{- range $i, $o := $o_array -}}
    {{- if eq $o.Camoflage "black"}}{{ if $printed_once }}, {{end}}"{{$o.Name}}"{{ $printed_once = true}}{{ end -}}
    {{- end -}}
];
private _{{$o_type}}_{{$o_mount}}_desert = [
    {{- $printed_once := false}}
    {{- range $i, $o := $o_array -}}
    {{- if eq $o.Camoflage "desert"}}{{ if $printed_once }}, {{end}}"{{$o.Name}}"{{ $printed_once = true}}{{ end -}}
    {{- end -}}
];
private _{{$o_type}}_{{$o_mount}}_woodland = [
    {{- $printed_once := false}}
    {{- range $i, $o := $o_array -}}
    {{- if eq $o.Camoflage "woodland"}}{{ if $printed_once }}, {{end}}"{{$o.Name}}"{{ $printed_once = true}}{{- end -}}
    {{- end -}}
];

private _{{$o_type}}_{{$o_mount}} = createHashMapFromArray [["nocamo", _{{$o_type}}_{{$o_mount}}_nocamo], ["black", _{{$o_type}}_{{$o_mount}}_black], ["desert", _{{$o_type}}_{{$o_mount}}_desert], ["woodland", _{{$o_type}}_{{$o_mount}}_woodland]];
_optics_{{ $o_type }} set ["{{ $o_mount }}", _{{$o_type}}_{{$o_mount}}];
{{- end }}

_arsenal_optics set ["{{ $o_type }}", _optics_{{ $o_type }}];
{{ end }}

_arsenal_optics