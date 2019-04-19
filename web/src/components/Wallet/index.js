import React from 'react';
import PerfectScrollbar from 'perfect-scrollbar';
import Selector from '../Selector';
import Tokens from './Tokens';
import Wrap from './Wrap';
import './styles.scss';

const OPTIONS = [
  { value: 'tokens', name: 'Tokens' },
  { value: 'wrap', name: 'Wrap' },
  { value: 'unwrap', name: 'Unwrap' }
];

class Wallet extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      selectedAccountID: OPTIONS[0].value
    };
  }

  render() {
    const { selectedAccountID } = this.state;
    return (
      <>
        <div className="title flex justify-content-between align-items-center">
          <div>Wallet</div>
          <Selector
            options={OPTIONS}
            selectedValue={selectedAccountID}
            handleClick={option => {
              this.setState({ selectedAccountID: option.value });
            }}
          />
        </div>
        <div className="flex-column flex-1 position-relative overflow-hidden" ref={ref => this.setRef(ref)}>
          {this.renderTabPanel()}
        </div>
      </>
    );
  }

  renderTabPanel() {
    const { selectedAccountID } = this.state;
    switch (selectedAccountID) {
      case 'tokens':
        return <Tokens />;
      case 'wrap':
        return <Wrap type="wrap" />;
      case 'unwrap':
        return <Wrap type="unwrap" />;
      default:
        return <Tokens />;
    }
  }

  setRef(ref) {
    if (ref) {
      this.ps = new PerfectScrollbar(ref, {
        suppressScrollX: true,
        maxScrollbarLength: 20
      });
    }
  }
}

export default Wallet;
