private _weapons_hash = createHashMap;

/*
 * 100 / (
 * RIS mount: /1.5;
 * Rareness: modern /2; specops used: /1.5; unusual hardware: /3.5;
 * Mounts: allows muzzle devs: /1.5; has preinstalled optics: /1.5; allows optics: /2;
 * front rail: /1,5; can mount grip/bipod: /1.5;
 * Integrated hardware: bipod: /2, int.toprail: /2, GL /2.5, int.front grip: /2;
 * Mechanics: balanced /4; high ROF: /2;
 * silenced: /1.5;
 * powerful cartridge: /2;
 * hi-cap mag: /2 (for pistols, precision rifles and shotguns)
 * )
*/
{{ $tpl_data := . }}
{{ range $w_type, $w_array := . }}
private _{{ $w_type | ToLower }}s = [
{{ range $i, $w := $w_array }}
{{- if eq $w_type "Starting" }}
    "{{ $w.Name }}"{{if IsLastItem $i (len $w_array)}}{{else}},{{end}}
{{- else }}
    ["{{ $w.Name }}", "{{ $w.MagType }}", "{{ $w.RailType }}", "{{ $w.Camoflage }}", {{ $w.AllowsMuzzleDevices }}, {{ $w.Features.AllowsGripsBipods }}],     {{if eq $w.RawScore 0.0}}{{ $w.Score }}{{else}}{{ $w.RawScore }}{{end}}{{if IsLastItem $i (len $w_array)}}{{else}},{{end}}
{{- end }}
{{- end }}
];
_weapons_hash set ["{{ $w_type | ToLower }}", _{{ $w_type | ToLower }}s];
{{ end }}

private _starting = [
{{- $printed_once := false}}
{{ range $w_type, $w_array := . }}
{{- range $i, $w := $w_array }}
{{- if eq $w.IsStarting 1 }}
    {{- if $printed_once }}, {{end}}["{{ $w.Name }}", "{{ $w.MagType }}", "{{ $w.RailType }}", "{{ $w.Camoflage }}", {{ $w.AllowsMuzzleDevices }}, {{ $w.Features.AllowsGripsBipods }}]{{ $printed_once = true}}
{{- end }}
{{- end }}
{{- end }}
];

[_weapons_hash, _starting]
