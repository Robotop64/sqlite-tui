# Views / Scripts

A **View** is understood in the context of this application as a user-defined layout of widgets, that enable the user to interact and view the contents of (data-)**Sources**.
These **Views** are defined in Lua scripts, representing a single **View** per script.

## Format

In principle the script may contain any valid Lua code, as well as certain predefined [variables](#variables) and [functions](#functions), that are recognized by the application.
Special is however that in the beginning of the script/file a `yaml`-like **header** is expected, giving meta-information about the script, such as its **Name** and **Type** and **Description**.
### Header
```yaml
---
Name: "<ViewName>"
Type: "View"
Description: "<ViewDescription>"
---
```

### Variables

The following variables are recognized in a **View** script and may be defined in the global scope to enable certain features of the application:
- `LoadSources` (table/list): A table/list of integers/indexes that declare which of the available **Sources** should be loaded for this **View**. 
- `Layout` (table): A table defining the layout of the **View**. See [Layouts](#layouts) for more information.

## Layouts

A **Layout** is a table declared in the global variable `Layout` akin to a [Widget](#widgets), that defines how the **View** and the contained **Widgets** are arranged.
All [Layouts](#layouts) contain the prefix letter `L` indicating that this element in the hierarchy is a **Layout**.
The following **Layouts** are currently supported:
- `LBox`: A simple box layout, arranging its children either vertically or horizontally.  
    Fields:
    - `dir`: Controlles if the children are arranged vertically (`vertical`) or horizontally (`horizontal`).
- `LBBox`: Short for Balanced Box. Similar to `LBox`, but the children maximize their size equally to fill the available space.
    Fields:
    - `dir`: Controlles if the children are arranged vertically (`vertical`) or horizontally (`horizontal`).
- `LFill`: A layout that fills the available space with a single child.
- `LWBox`: A weighted box layout, arranging its children either vertically or horizontally, where each child may be assigned a weight to control how much space it takes relative to its siblings.
    Fields:
    - `dir`: Controlles if the children are arranged vertically (`vertical`) or horizontally (`horizontal`).
    - `weights`: A table/list of integers defining the weights for each child. The length of this table must match the number of children. There are three kinds of weights with priority in the following order:
        - `-1`: The child according to this weight is trying to take as little space as possible.
        - `>0`: The child according to this weight will take space relative to the sum of all weights.
        - `0`: The child according to this weight will try to take as much space as possible.

All **Layouts** only restrict/manage the along their primary axis (defined by `dir`), while the secondary axis is always maximized to fill the available space.
To achieve a more complex layout, **Layouts** may be nested within each other.

For all **Layouts** elements of lower hierarchy are defined as an anonymous table in the lastfield.
### Example Layout

```lua
LWBox {
    dir = "horizontal",
    weights = {0, -1},
    {
        OtherLayout {
            ...
        },
        Widget {
            ...
        },        
    }
}
```

## Widgets

This section describes a series of widgets a user may define in a script to customize a **View**.

### Table

The Table widget is declared as a LuaTable with the following possible/recognized fields:
```lua
WTable {
    source = 0,
    enable = {
        title = {
            text = "<TableName>",
            alignment = "<left|center|right>",
            style = {
                bold = true,
                italic = false,
            },
        },
        headers = {
            column = true,
            row    = true,
        },
        filter = {
            table = true,
            columns = true,
            sort_by = true,
            sort_dir = true,
            group_by = true,
            filter = true,
            limit = true,
            page = true,
        },
    },
    editable = true,
}
```