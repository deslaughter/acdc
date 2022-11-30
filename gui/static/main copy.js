export default {
    el: "#app",
    data() {
        return {
            count: 0,
            schema: [],
            inp: {
                Defaults: {},
            },
            boolOptions: [
                { text: 'True', value: true },
                { text: 'False', value: false },
            ],
            loaded:false,
        }
    },
    mounted() {
        axios.get('/fasted/api/schemas/Fast').then(response => {
            let groups = [];
            let group = {};
            for (const e of response.data) {
                if ("Options" in e) {
                    for (const opt of e.Options) {
                        if (opt.text != opt.value) {
                            opt.text = opt.value + ": " + opt.text;
                        }
                    }
                }
                if ("Heading" in e) {
                    group = { name: e.Heading, entries: [] };
                    groups.push(group);
                } else {
                    group.entries.push(e);
                }
            }
            this.schema = groups;
            console.log(response);
        }).catch(error => {
            console.log(error);
        })
    },
    methods: {
        toggleDefault(entry) {
            if (entry.Keyword in this.inp.Defaults) {
                this.$delete(this.inp.Defaults, entry.Keyword);
            } else {
                this.$set(this.inp.Defaults, entry.Keyword, {});
            }
        },
        entryValid(entry) {
            if (entry.Keyword in this.inp.Defaults) return true;
            const value = this.inp[entry.Field];
            if (value == null) return false;
            switch (entry.Type) {
                case 'bool':
                    return value == true || value == false;
                case 'string':
                    return value.length > 0;
                case 'int':
                    return Number.isInteger(value);
                case 'float64':
                    return Number.isFinite(value);
            }
            return true;
        },
        entryActive(entry) {
            if ("Active" in entry) {
                for (const c of entry.Active) {
                    switch (c.relation) {
                        case "==":
                            if (this.inp[c.field] == c.value) continue;
                            break;
                        case "!=":
                            if (this.inp[c.field] != c.value) continue;
                            break;
                        case "<":
                            if (this.inp[c.field] < c.value) continue;
                            break;
                        case ">":
                            if (this.inp[c.field] > c.value) continue;
                            break;
                        case "in":
                            if (c.value.includes(this.inp[c.field])) continue;
                            break;
                    }
                    return false;
                }
            }
            return true;
        }
    }
}