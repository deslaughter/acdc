export default {
    el: "#app",
    data() {
        return {
            analysis: {},
            modelFiles: [],
            loaded: false,
            conditions: {
                WindSpeed: 0,
                RotorSpeed: 0,
                BladePitch: 0,
                TowerTopDispForeAft: 0,
                TowerTopDispSideSide: 0,
            }
        }
    },
    mounted() {
        axios.get('/acdc/api/analysis').then(response => {
            this.analysis = Object.assign({}, this.analysis, response.data)
            console.log(response);
            this.loaded = true;
        }).catch(error => {
            console.log(error);
        })
    },
    methods: {
        filterModelFiles(e) {
            this.modelFiles = Array.from(e.target.files).filter(file =>
                !file.name.startsWith('.'));
            this.modelFiles.sort((a, b) => (a.$path > b.$path) ? 1 : -1);
        },
        importModel() {
            let formData = new FormData();
            for (let file of this.modelFiles) {
                formData.append("paths", file.$path);
                formData.append("files", file, file.$path);
            }
            axios.post('/acdc/api/model', formData, {
                headers: {
                    "Content-Type": "multipart/form-data",
                }
            }).then(response => {
                this.modelFiles = [];
                this.analysis.Turbine = response.data;
                console.log(response)
            }).catch(error => {
                console.log(error);
            })
        },
        addConditions() {
            this.analysis.Conditions.push(this.conditions);
            this.updateConditions();
        },
        removeConditions(i) {
            this.analysis.Conditions.splice(i, 1);
            this.updateConditions();
        },
        updateConditions(conditions) {
            axios.post('/acdc/api/conditions', this.analysis.Conditions).then(response => {
                this.analysis.Conditions = response.data;
                console.log(response)
            }).catch(error => {
                console.log(error);
            })
        },
    }
}