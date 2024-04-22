# Gourlex

One domain:
```
gourlex -t https://github.com -s -uO > gourlex_results
nuclei -l gourlex_results -tags token,tokens,takeover -es unknown -rl 500 -c 100 -silent
```
Multiple domains (via file)
```
gourlex -f domains.txt -s -uO > gourlex_results
nuclei -l gourlex_results -tags token,tokens,takeover -es unknown -rl 500 -c 100 -silent
```
![image](https://github.com/reewardius/gourlex/assets/68978608/e4123163-368c-498b-b0c6-00ec999da068)

## Overview

Gourlex is a simple tool that can be used to extract URLs and paths from web pages. 
It can be helpful during web application assessments to uncover additional targets.

![gourlex](https://github.com/trap-bytes/gourlex/blob/main/static/gourlex.png)

## Features

- **URLs and Paths Extraction**
  - The tool can be used to extract only URLs, only paths, or both.
- **Silent mode for easy integration with other tools**
  - The tool provides a silent mode, making it easy to integrate its output into other tools during the reconnaissance and enumeration phases.

## Install

```
go install github.com/trap-bytes/gourlex@latest
```
## Usage:

```
gourlex -h
```

This will display help for the tool. Here are all the arguments it supports.

```
Usage:
  gourlex [arguments]

The arguments are:
  -t string    Specify the target URL (e.g., domain.com or https://domain.com)  
  -p string    Specify the proxy URL (e.g., 127.0.0.1:8080)
  -c string    Specify cookies (e.g., user_token=g3p21ip21h; 
  -r string    Specify headers (e.g., Myheader: test
  -s           Silent Mode, avoid printing banner and other messages
  -uO          Extract only full URLs
  -pO          Extract only URL paths
  -h           Display help

Example:
  gourlex -t domain.com
```
# Nuclei (custom)
```
gourlex -t github.com -s -uO > gourlex_results
nuclei -l gourlex_results -tags token,tokens,takeover,provider -es unknown -rl 500 -c 100 -silent
```
# gourlex + nuclei (windows)
```
Get-Content domains | ForEach-Object { gourlex -t $_ -s -uO | Out-File -Append gourlex_results }
nuclei -l gourlex_results -tags token,tokens,takeover,provider -es unknown -rl 500 -c 100 -silent
```
# gourlex + nuclei (linux)
```
while IFS= read -r line; do gourlex -t $line -s -uO; done < domains >> gourlex_results
nuclei -l gourlex_results -tags token,tokens,takeover,provider -es unknown -rl 500 -c 100 -silent
```
# gourlex + nuclei (firebase)
Взято из статей по взлому жепы Firebase. Это всё один ресёрч, разбитый по трём блогам 
1. https://kibty.town/blog/chattr/
2. https://mrbruh.com/chattr/
3. https://env.fail/posts/firewreck-1/
```
nuclei -l gourlex_results -id firebase-config-exposure -rl 500 -c 100 -silent
```
