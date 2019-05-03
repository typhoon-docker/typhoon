const classNameRegex = /^\$class-/;

const parseClassNames = props => {
  const classNames = [];
  const newProps = {};
  Object.keys(props).forEach(propName => {
    if (propName.match(classNameRegex)) {
      if (props[propName]) {
        classNames.push(propName.substr(7));
      }
    } else {
      newProps[propName] = props[propName];
    }
  });
  return {
    classNames,
    props: newProps,
  };
};

export default parseClassNames;
