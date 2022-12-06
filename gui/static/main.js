

export default {
    el: "#app",
    data() {
        return {
            analysis: {},
            modelFiles: [],
            loaded: false,
            modelPathValid: false,
            execPathValid: false,
            ws: {},
            evalStatus: [],
            conditions: {
                WindSpeed: 0,
                RotorSpeed: 0,
                BladePitch: 0,
                TowerTopDispForeAft: 0,
                TowerTopDispSideSide: 0,
            },
            importAlert: '',
            cpuOpts: Array.from({ length: 24 }, (v, k) => k + 1),
            trueFalseOpts: [
                { value: true, text: 'True' },
                { value: false, text: 'False' },
            ],
            trimCaseOpts: [
                { value: 1, text: '1 - Yaw' },
                { value: 2, text: '2 - Torque' },
                { value: 3, text: '3 - Pitch' },
            ],
            linInputsOpts: [
                { value: 0, text: '0 - None' },
                { value: 1, text: '1 - Standard' },
                { value: 2, text: '2 - All module inputs (debug)' },
            ],
            linOutputsOpts: [
                { value: 0, text: '0 - None' },
                { value: 1, text: '1 - From OutList(s)' },
                { value: 2, text: '2 - All module outputs (debug)' },
            ],
            compElastOpts: [
                { value: 1, text: '1 - ElastoDyn' },
                { value: 2, text: '2 - ElastoDyn + BeamDyn for blades' },
            ],
            compAeroOpts: [
                { value: 0, text: '0 - None' },
                { value: 1, text: '1 - AeroDyn v14' },
                { value: 2, text: '2 - AeroDyn v15' },
            ],
            compServoOpts: [
                { value: 0, text: '0 - None' },
                { value: 1, text: '1 - ServoDyn' },
            ],
            compHydroOpts: [
                { value: 0, text: '0 - None' },
                { value: 1, text: '1 - HydroDyn' },
            ],
            compSubOpts: [
                { value: 0, text: '0 - None' },
                { value: 1, text: '1 - SubDyn' },
                { value: 2, text: '2 - External Platform MCKF' },
            ],
            interpOrderOpts: [
                { value: 1, text: '1 - Linear' },
                { value: 2, text: '2 - Quadratic' },
            ],
            outFileFmtOpts: [
                { value: 1, text: '1 - Text [<RootName>.out]' },
                { value: 2, text: '2 - Binary [<RootName>.outb]' },
                { value: 3, text: '3 - Both' },
            ],
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
        const wsURL = "ws://" + document.location.host + "/acdc/api/evaluate";
        this.ws = new ReconnectingWebSocket(wsURL);
        this.ws.onmessage = (evt) => {
            console.log(evt.data)
            this.evalStatus = JSON.parse(evt.data)
        };
    },
    methods: {
        debounceUpdateAnalysis: _.debounce(function () {
            this.updateAnalysis();
        }, 1000),
        updateAnalysis() {
            axios.put('/acdc/api/analysis', this.analysis).then(response => {
                this.analysis = response.data;
                console.log(response.data)
            }).catch(error => {
                console.log(error);
            })
        },
        importModel() {
            let formData = new FormData();
            formData.append("path", this.analysis.ModelPath);
            axios.post('/acdc/api/model', formData, {
                headers: { "Content-Type": "multipart/form-data" }
            }).then(response => {
                this.analysis.Model = response.data;
                this.importAlert = '';
                console.log(response.data)
            }).catch(error => {
                this.importAlert = error.response.data
                console.log(error);
            })
        },
        addConditions() {
            this.analysis.Conditions.push(this.conditions);
            this.updateAnalysis();
        },
        removeConditions(i) {
            this.analysis.Conditions.splice(i, 1);
            this.updateAnalysis();
        },
        evalStart() {
            axios.post('/acdc/api/evaluate', this.analysis).then(response => {
                console.log(response)
            }).catch(error => {
                console.log(error);
            })
        },
        evalCancel() {
            axios.delete('/acdc/api/evaluate').then(response => {
                console.log(response)
            }).catch(error => {
                console.log(error);
            })
        },
        validateModelPath() {
            let formData = new FormData();
            formData.append("path", this.analysis.ModelPath);
            const resp = axios.post('/acdc/api/validate-path', formData, {
                headers: { "Content-Type": "multipart/form-data" }
            }).then(response => {
                this.modelPathValid = true;
            }).catch(error => {
                this.modelPathValid = false;
            })
        },
        validateExecPath() {
            let formData = new FormData();
            formData.append("path", this.analysis.ExecPath);
            const resp = axios.post('/acdc/api/validate-path', formData, {
                headers: { "Content-Type": "multipart/form-data" }
            }).then(response => {
                this.execPathValid = true;
            }).catch(error => {
                this.execPathValid = false;
            })
        },
        updateLinTimes(NLinTimes) {
            if (!this.analysis.Model.FAST.CalcSteady) {
                this.analysis.Model.FAST.LinTimes.length = NLinTimes;
            }
        },
    }
}