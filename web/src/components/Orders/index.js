import React from 'react';
import OpenOrders from './OpenOrders';
import Trades from './Trades';
import Selector from '../Selector';

const OPTIONS = [{ value: 'openOrders', name: 'Open' }, { value: 'filled', name: 'Filled' }];

class Orders extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      selectedType: OPTIONS[0].value
    };
  }
  render() {
    const { selectedType } = this.state;
    return (
      <>
        <div className="title flex justify-content-between align-items-center">
          <div>
            <div>Orders</div>
            <div className="text-secondary">View your open orders</div>
          </div>
          <Selector
            options={OPTIONS}
            selectedValue={selectedType}
            handleClick={option => {
              this.setState({ selectedType: option.value });
            }}
          />
        </div>
        {selectedType === 'openOrders' ? <OpenOrders /> : <Trades />}
      </>
    );
  }
}

export default Orders;
