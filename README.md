# csv-dl

`csv-dl` reads CSV on stdin, and downloads linked files.

<!-- TOC -->

- [Get it](#get-it)
- [Use it](#use-it)

<!-- /TOC -->

## Get it

Using go get:

```bash
go get -u github.com/sgreben/csv-dl
```

Or [download the binary](https://github.com/sgreben/csv-dl/releases/latest) from the releases page.

```bash
# Linux
curl -LO https://github.com/sgreben/csv-dl/releases/download/0.0.1/csv-dl_0.0.1_linux_x86_64.zip
unzip csv-dl_0.0.1_linux_x86_64.zip

# OS X
curl -LO https://github.com/sgreben/csv-dl/releases/download/0.0.1/csv-dl_0.0.1_osx_x86_64.zip
unzip csv-dl_0.0.1_osx_x86_64.zip

# Windows
curl -LO https://github.com/sgreben/csv-dl/releases/download/0.0.1/csv-dl_0.0.1_windows_x86_64.zip
unzip csv-dl_0.0.1_windows_x86_64.zip
```

## Use it

```text
Usage of csv-dl:
  -f	(alias for -force-overwrite)
  -force-overwrite
    	overwrite existing files
  -l value
    	(alias for -link)
  -link value
    	a link to download, may use go {{template}} syntax and refer to data columns by index (column i) or name (field "f")
  -p int
    	(alias for -parallel) (default 8)
  -parallel int
    	number of parallel connections (default 8)
  -q	(alias for -quiet)
  -quiet
    	suppress all logging
  -s string
    	(alias for -schema)
  -schema string
    	use the given CSV expression as the table schema
  -skip-csv-eader
    	assume the first row is the CSV header, skip it
  -u	(alias for -use-csv-header)
  -use-csv-header
    	assume the first row is the CSV header, use it as a schema
```
