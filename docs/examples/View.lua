---
Name: B
Type: "View"
Description: "Widget Test"
---

load_sources = {0}

query = {
    master = "SELECT * FROM my_table"
}

layout = LWBox{
    dir = "vertical",
    weights = {0.0},
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
                    link = "0",
                },
                WSeparator{},
                LWBox{
                    dir = "vertical",
                    weights = {0.0},
                    {
                        WFilter{
                            type = "table",
                            enable = {
                                table = false,
                                columns = true,
                                sort_by = true,
                                sort_dir = true,
                                group_by = true,
                                filter = true,
                                limit = true,
                                page = true,
                                query = false,
                            },
                            link = "0",
                        },
                    }
                },
            },
        }
    }
}
