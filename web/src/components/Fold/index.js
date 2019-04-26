import React from 'react';
import UnderlineTab from '../UnderlineTab';
import './styles.scss';

export default class Fold extends React.PureComponent {
  constructor(props) {
    super(props);

    this.state = {
      selectedIndex: 0
    };
  }
  render() {
    const children = this.props.children;
    const child = children[this.state.selectedIndex];
    const options = [];

    for (let i = 0; i < children.length; i++) {
      const child = children[i];
      options.push({
        title: child.props['data-fold-item-title'],
        onClick: () => {
          this.setState({
            selectedIndex: i
          });
        }
      });
    }

    return (
      <div className={[this.props.className, 'fold'].join(' ')}>
        <div className="flod-header">
          <span>{this.props.title}</span>
          <UnderlineTab options={options} selectedIndex={this.state.selectedIndex} />
        </div>
        {child}
      </div>
    );
  }
}
