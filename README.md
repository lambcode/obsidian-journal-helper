# Daily Journal Script

This Go program will create a new journal markdown file for each day that it is run. 
It copies over incomplete tasks from the most recent day. Each time a task is copied, an
asterisk will be added to the task easily see how many times it has been deferred.

## Building

```shell
go build .
```

## Running

```shell
export JOURNAL_DIR="<some directory here>"
./journal
```