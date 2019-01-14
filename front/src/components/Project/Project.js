import React, { useRef } from "react";

const Project = ({ project: { id, name, repository_url, template_id }, onSelect, selected }) => (
  <details key={id} open={selected}>
    <summary onClick={onSelect}>{name}</summary>
    <p>
      Projet en {template_id}, accessible via <a href={repository_url}>{repository_url}</a>
    </p>
  </details>
);

export default Project;
