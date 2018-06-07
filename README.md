# check-elasticsearch
Icinga Checks based on ElasticSearch results

## Default Flags

| flag | short | description |
| -- | -- | -- |
| critical | c | defines a [threshold](#threshold) for a critical return status |
| warning | w | defines a [threshold](#threshold) for a critical return status |
| index | i | the ElasticSearch index to search in |
| date-range | d | the date range to filter results |

## Thresholds

<https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT>

| Range definition |	Generate an alert if x... |
| -- | -- |
| 10 | < 0 or > 10, (outside the range of {0 .. 10}) |
| 10: | < 10, (outside {10 .. ∞}) |
| ~:10 | > 10, (outside the range of {-∞ .. 10}) |
| 10:20 | < 10 or > 20, (outside the range of {10 .. 20}) |
| @10:20 | ≥ 10 and ≤ 20, (inside the range of {10 .. 20} |