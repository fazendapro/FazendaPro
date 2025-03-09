defmodule Farmnotifications.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      FarmnotificationsWeb.Telemetry,
      {DNSCluster, query: Application.get_env(:farmnotifications, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: Farmnotifications.PubSub},
      # Start the Finch HTTP client for sending emails
      {Finch, name: Farmnotifications.Finch},
      # Start a worker by calling: Farmnotifications.Worker.start_link(arg)
      # {Farmnotifications.Worker, arg},
      # Start to serve requests, typically the last entry
      FarmnotificationsWeb.Endpoint
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Farmnotifications.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    FarmnotificationsWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
