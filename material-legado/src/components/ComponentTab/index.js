import React, { Component } from 'react';
import { AppConsumer } from '../../service/AppContext';
import { Container, Item } from './styles';

export default class ComponentTab extends Component {

    async onDragStart(event, component) {
        event.dataTransfer.effectAllowed = 'copy';
        event.dataTransfer.dropEffect = 'copy';
        event.dataTransfer.setData("text/plain", event.target.dataset.id);

        if (component.preview && component.preview.image) {
            event.dataTransfer.setDragImage(component.preview.image, 0, 0);
        }
    }

    render() {
        return <AppConsumer>
            {
                data => (
                    <ul style={Container}>
                        {
                            data.components.map((data) => (
                                <li style={Item} key={data.id} draggable="true" data-id={data.id} onDragStart={(event) => this.onDragStart(event, data)}>{data.name}</li>
                            ))
                        }
                    </ul>
                )
            }
        </AppConsumer>;
    }
}