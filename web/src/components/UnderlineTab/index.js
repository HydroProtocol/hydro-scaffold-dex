import React from 'react';
import './styles.scss';

export default class UnderlineTab extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      underlineClassName: 'underline'
    };
    this.mounted = false;
  }

  componentDidMount() {
    window.addEventListener('load', () => {
      this.forceUpdate();
    });
    this.mounted = true;
  }

  componentWillUnmount() {
    this.mounted = false;
  }

  render() {
    const { options, selectedIndex } = this.props;
    const optionsElements = [];

    for (let i = 0; i < options.length; i++) {
      const option = options[i];
      optionsElements.push(
        <div key={option.title} className={`item${i === selectedIndex ? ` active` : ''}`} onClick={option.onClick}>
          {option.title}
          {!this.container && i === selectedIndex ? <div className="defaultUnderline" /> : null}
        </div>
      );
    }

    let left, width;

    if (this.container) {
      const activeItem = this.container.children[selectedIndex];
      left = activeItem.offsetLeft;
      width = activeItem.offsetWidth;
    } else {
      left = 0;
      width = 0;
    }

    const underline = <div style={{ left, width }} className={this.state.underlineClassName} />;

    return (
      <div className="underlineTabContainer" ref={this.ref}>
        {optionsElements}
        {underline}
      </div>
    );
  }

  ref = ref => {
    if (!ref) {
      return;
    }
    this.container = ref;
    setTimeout(() => {
      if (this.mounted) {
        this.setState({
          underlineClassName: 'underline transition'
        });
      }
    }, 300);
  };
}
