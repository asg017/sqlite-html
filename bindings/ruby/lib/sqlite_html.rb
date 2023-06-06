require "version"

module SqliteHtml
  class Error < StandardError; end
  def self.html_loadable_path
    File.expand_path('../html0', __FILE__)
  end
  def self.load(db)
    db.load_extension(self.html_loadable_path)
  end
end
