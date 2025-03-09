defmodule FarmNotificationsWeb.PageController do
  use FarmNotificationsWeb.WhatsApp

  def create(conn, %{"event" => "novo_gado", "data" => data}) do
    FarmNotificationsWeb.Endpoint.broadcast("gado:notifications", "novo_gado", data)
    json(conn, %{status: "ok"})
  end
end
