-- ---
-- Name: B
-- Type: "View"
-- Description: "Widget Test"
-- ---

LoadSources = { 0 }

query = {
    master = "SELECT * FROM my_table"
}

Layout = LWBox {
    dir = "vertical",
    weights = { 0.0 },
    {
        LWBox {
            dir = "horizontal",
            weights = { 0.0, -1.0, -1.0 },
            {
                WTable {
                    source = 0,
                    enable = {
                        title = {
                            text = "My Table:",
                            alignment = "left",
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
                },
            },
        }
    }
}
