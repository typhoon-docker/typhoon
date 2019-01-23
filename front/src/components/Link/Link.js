import React from 'react';

import Underline from '../Underline';

const Link = ({ newTab, ...props }) => {
  const target = newTab ? '_blank' : '_self';

  return <Underline {...props} variant="dashed" as="a" target={target} />;
};

export default Link;
