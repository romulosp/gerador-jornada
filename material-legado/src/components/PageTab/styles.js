const Container = {
    listStyle: 'none',
    margin: 0,
    padding: 0,
    marginTop: 20,
    borderTop: '1px solid #545454',
    overflow: 'auto',
    maxHeight: 'calc(100vh - 167px)'
};
const Item = {
    color: 'rgb(186, 189, 198)',
    fontSize: 14,
    textTransform: 'capitalize',
    fontWeight: 'bolder',
    backgroundColor: 'rgb(36, 41, 49)',
    padding: '10px 20px',
    borderBottom: '1px solid #545454'
};
const ButtonContainer = {
    display: 'flex',
    justifyContent: 'flex-end'
};

const Button = {
    backgroundColor: 'rgb(36, 41, 49)',
    marginRight: 15,
    border: 'none',
    color: '#BABDC6',
    height: 30,
    fontSize: 24
};

export {
    Container,
    Item,
    ButtonContainer,
    Button
}