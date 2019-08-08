import React from 'react';
const Bar = ({ percent }) => {
  return (
    <div className="bar" style={{ minWidth: `${percent}%` }} ></div>
  )
}

export default Bar;