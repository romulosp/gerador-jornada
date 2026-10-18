import React, { Component } from 'react';
import { Container } from './styles';
export default class PageElementPreview extends Component {

    convertToArray = (data) => {
        if (!data) {
            return []
        }
        return Object.keys(data).map(key => data[key]);
    };

    render() {
        const { element, data } = this.props;
        const properties = element.properties;
        const renderElement = element.render(properties, data);
        return (
            <div style={{ ...Container, ...element.getContainerStyle()}}>
                {
                    renderElement
                }
            </div>
        );
    }
}