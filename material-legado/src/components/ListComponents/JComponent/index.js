import uuidv4 from'uuid/v4';

export default class JComponent {
    constructor(name,props) {
        this.name = name;
        this.id = uuidv4();

        this.properties = {
            width: 20,
            height: 20,
            top: props.y + 10,
            left: props.x + 10,
            name: JComponent + this.id,
            description: '',
            events: false
        };

        this.schema = {
            basic: {
                label: 'Básico',
                data: {
                    name: {
                        label: 'Nome',
                        key: 'name',
                        type: 'string'
                    },
                    description: {
                        label: 'Descrição',
                        key: 'description',
                        type: 'string'
                    },
                }
            },
            size: {
                label: 'Tamanho',
                data: {
                    width: {
                        label: 'Largura',
                        key: 'width',
                        type: 'number',
                        isStyle: true,
                        renderStyle: (val) => `${val}px`,
                        min: 1,
                        max: 800
                    },
                    height: {
                        label: 'Altura',
                        key: 'height',
                        type: 'number',
                        isStyle: true,
                        renderStyle: (val) => `${val}px`,
                        min: 1,
                        max: 800
                    }
                }
            },
            position: {
                label: 'Posição',
                data: {
                    top: {
                        label: 'Vertical',
                        key: 'top',
                        isContainerStyle: true,
                        renderStyle: (val) => `${val}px`,
                        type: 'number',
                        min: 0,
                        max: 600
                    },
                    left: {
                        label: 'Horizontal',
                        key: 'left',
                        isContainerStyle: true,
                        renderStyle: (val) => `${val}px`,
                        type: 'number',
                        min: 0,
                        max: 800
                    }
                }
            }
        };
    }

    toJSON(){
        return {name: this.name, id: this.id, properties: this.properties}
    }
    
    fromJSON(){
        
    }

    _convertToArray = (data) => {
        if (!data) {
            return []
        }
        return Object.keys(data).map(key => data[key]);
    };

    getSchemaElements = () => {
        const schemaFolder = this._convertToArray(this.schema);
        const elements = schemaFolder.map(item => this._convertToArray(item.data));
        return elements.reduce((acc, val) => acc.concat(val), []);
    }

    _getStyle = (fn) => {
        const elements = this.getSchemaElements().filter(fn);
        return elements.map(item => {
            return {
                key: item.key,
                value: item.renderStyle(this.properties[item.key])
            };
        }).reduce((obj, item) => {
            obj[item.key] = item.value
            return obj
        }, {});
    }

    getContainerStyle = () => {
        return this._getStyle(item => item.isContainerStyle);
    }

    getElementStyle = () => {
        return this._getStyle(item => item.isStyle);
    }
}
