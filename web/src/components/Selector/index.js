import React from 'react';
import './styles.scss';

class Selector extends React.PureComponent {
  render() {
    const { options, selectedValue, handleClick } = this.props;
    if (!options) {
      return null;
    }
    return (
      <div className="selector">
        <ul className="nav nav-tabs">
          {options.map(option => {
            return (
              <li
                key={option.value}
                className={`nav-item${selectedValue === option.value ? ' active' : ''}`}
                onClick={() => handleClick(option)}>
                <div className="text-center">{option.name}</div>
              </li>
            );
          })}
        </ul>
      </div>
    );
  }
}

export default Selector;
