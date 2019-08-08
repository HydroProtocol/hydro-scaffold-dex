import React from 'react';
import './styles.scss';
import Bar from './Bar'

function BarChart ({ percent }){
  return <Bar percent={percent}/>;
}
export default BarChart;