# flogo-jq-activity
Filter json using jq with a proper multiline textbox in the GUI for editing the query, like the JSexec activity

In the descriptor.json, a schema is used for the GUI, which even though it says draft v4, is not a draft v4 schema, but something a developer came up with to allow drop-downs for types in the UI.

The schema can be found in the file [transformation/activity/jq/argumentnames.json](transformation/activity/jq/argumentnames.json).