# check-elasticsearch
Icinga Checks based on ElasticSearch results

## Commands

### stringQuery

Submit any Lucene query to search for results in the given index.

- Search for any message cntaining the string "ldap" within the last 15 minutes
```
message:("ldap") AND @timestamp:>now-15m
```

## Default Flags

| Flag | Short | Description |
| -- | -- | -- |
| index | i | The ElasticSearch index to search in |
| critical | c | Defines a [threshold](#thresholds) for a critical return status |
| warning | w | Defines a [threshold](#thresholds) for a critical return status |
| debug | d | Switch debug mode on/off |

## Thresholds

<https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT>

| Range definition |	Generate an alert if x... |
| -- | -- |
| 10 | < 0 or > 10, (outside the range of {0 .. 10}) |
| 10: | < 10, (outside {10 .. ∞}) |
| ~:10 | > 10, (outside the range of {-∞ .. 10}) |
| 10:20 | < 10 or > 20, (outside the range of {10 .. 20}) |
| @10:20 | ≥ 10 and ≤ 20, (inside the range of {10 .. 20} |