private _vehicles = createHashMap;

{{ range $role, $veh_array := . }}

private _{{ $role }} = [];

{{- range $i, $veh := $veh_array }}
private _veh_{{$i}} = createHashMapFromArray [
    ["type", "{{ $veh.Type }}"], ["class", "{{ $veh.Class }}"],
    ["family", "{{ $veh.Family }}"], ["name", "{{ $veh.Name }}"],
    ["cost", {{ $veh.Cost }}], ["fuel", {{ $veh.FuelCost }}], ["upgrade_cost", {{ $veh.UpgradeCost }}],
    ["woodland", {{ $veh.CamosWoodland.CamoList }}], ["desert", {{ $veh.CamosDesert.CamoList }}], ["nocamo", {{ $veh.CamosNocamo.CamoList }}],
    ["unlocked", {{ $veh.IsUnlocked }}]{{ if $veh.RawProperties }}, {{ $veh.RawProperties }}{{ end }}
];
_{{ $role }} append [_veh_{{$i}}, {{ $veh.Rareness }}];
{{- end }}
_vehicles set ["{{ $role }}", _{{ $role }}];
{{- end }}

_vehicles
