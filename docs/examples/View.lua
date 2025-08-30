---
Name: B
Type: "View"
Description: "Widget Test"
---

sources = {0}

layout = LWBox{
    dir = "vertical",
    weights = {0.75},
    {
        LWBox{
            dir = "horizontal",
            weights = {0.0,-1.0,-1.0},
            {
                WTable{
                    idx_source = 0,
                    title = {
                        text = "My Table",
                        alignment = "left",
                        style = {bold = true},
                    },
                    header = {
                        column = true,
                        row = false,
                    },
                    editable = true,
                    link = "0"
                },
                WSeparator{},
                LWBox{
                    dir = "vertical",
                    weights = {-1.0,0.0},
                    {
                        WFilter{
                            type = "table",
                            enable = {
                                tables = false,
                                columns = true,
                                filter_by = true,
                                limit = true,
                                direction = true,
                                page = true,
                            }
                            link = "0"
                        },
                    }
                },
            },
        }
    }
}
