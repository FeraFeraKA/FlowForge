export function PageHeader() {
  return (
    <header className="page-header">
      <div>
        <div className="eyebrow">Pipelines / Web application</div>
        <h1>Build and deploy</h1>
        <p>Install, test, build and safely deliver the web application.</p>
      </div>
      <div className="page-header__actions">
        <button className="button button--secondary" type="button">Edit pipeline</button>
        <button className="button button--primary" type="button">
          <span aria-hidden="true">▶</span> Run pipeline
        </button>
      </div>
    </header>
  )
}
