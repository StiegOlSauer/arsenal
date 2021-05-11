# arsenal - Complementary arsenal generator for Liberation mission for Arma 3

Generates SQF code from provided CSV tables which is then consumed by Arma 3 Liberation mission (https://github.com/StiegOlSauer/greuh_liberation.Altis).

Usage example:
1. Get .ods equipment table from Liberation repository
2. Make required edits and export .csv files into some convenient place (i.e. `arsenal/pmcafrf/weapons.csv`)
3. Run the generator from that directory:

  `$ ../arsenal -w --weapons-csv weapons.csv --templates-dir /path/to/tpl/`

4. SQF files in corresponding subdirs will be created in `.` directory.
5. Copy the directory into Liberation mission `scripts/gameplay_templates/arsenal/`
