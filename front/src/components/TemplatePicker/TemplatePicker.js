import React, { Fragment } from 'react';

import templates from '/utils/templates.json';

import node from './node.png';
import php from './php.png';
import python3 from './python3.svg';
import raw from './static.svg';
import cra from './create-react-app.svg';

import { img, input, label, wrapper, name } from './TemplatePicker.css';

const images = {
  node,
  php,
  python3,
  static: raw,
  'create-react-app': cra,
};

const TemplatePicker = ({ onSelect }) => {
  return (
    <>
      <h2>En quoi as-tu codé ton projet ?</h2>
      <div className={wrapper}>
        {Object.values(templates).map(template => (
          <Fragment key={template.template_id}>
            <input
              type="radio"
              id={template.template_id}
              name="template_id"
              value={template.template_id}
              className={input}
              onChange={() => onSelect(template)}
            />
            <label htmlFor={template.template_id} className={label}>
              <img
                src={images[template.template_id]}
                alt={template.template_id}
                className={img}
                title={template.template_id}
              />
              <span className={name}>{template.template_id}</span>
            </label>
          </Fragment>
        ))}
      </div>
    </>
  );
};

export default TemplatePicker;
