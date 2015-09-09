# fluentd forward plugin for SensorBee

fluentd plugin support receiving tuples from fluentd out_forward plugin and
sending tuples to in_forward plugin.

## Registering plugin

Just import plugin package from an application:

```go
import (
    _ "pfi/sensorbee/fluentd/plugin"
)
```

Or, manually register the source and the sink to bql package.

## Creating a source

```
CREATE SOURCE <name> TYPE fluentd;
```

### Optional parameters

#### bind (string)

`bind` parameter has the address and the port number on which the sink listens.
Its format is `<addr>:<port>`. For example `127.0.0.1:24224`, `0.0.0.0:12345`.

Its default value is `'127.0.0.1:24224'`.

Example:

```
CREATE SOURCE metrics TYPE fluentd WITH bind='0.0.0.0:8080';
```

#### tag_field (string)

`tag_field` parameter has a name of a field in a tuple which contains the tag
added by fluentd.

Its default value is `'tag'`.

Example:

```
CREATE SOURCE metrics TYPE fluentd WITH tag_field='fluentd_tag';
```

By setting this, `tuple.fluentd_tag` will have the value like `'system.test'`.

## Creating a sink

```
CREATE SINK <name> TYPE fluentd;
```

### Optional parameters

#### tag_field (string)

`tag_field` parameter has a name of a field in a tuple which contains the tag
added by fluentd. If a tuple contains the field, the sink use its value as a
tag. If a tuple doesn't have the field or its value is empty, the default tag
will be used instaed.

Its default value is `'tag'`.

Example:

```
CREATE SINK out_forward TYPE fluentd WITH tag_field='fluentd_tag';
```

#### default_tag (string)

`default_tag` parameters has a default tag name which is used when a tuple
doesn't have tag field or the field is empty.

Its default value is `'sensorbee'`.

Example:

```
CREATE SINK out_forward TYPE fluentd WITH default_tag='sensorbee.forward';
```

TODO: Add other parameters
