defmodule FarmNotificationsWeb.Services.WhatsApp do
  require Logger

  @whatsapp_api_version "v20.0"
  @whatsapp_base_url "https://graph.facebook.com/#{@whatsapp_api_version}"

  @spec send_message(String.t(), String.t()) :: {:ok, map()} | {:error, atom() | String.t()}
  def send_message(phone_number, message) when is_binary(phone_number) and is_binary(message) do
    with {:ok, phone_number} <- validate_phone_number(phone_number),
         {:ok, message} <- validate_message(message),
         {:ok, config} <- get_config() do
      do_send_message(phone_number, message, config)
    end
  end

  defp do_send_message(phone_number, message, config) do
    url = "#{@whatsapp_base_url}/#{config.phone_number_id}/messages"

    body = %{
      messaging_product: "whatsapp",
      to: phone_number,
      type: "text",
      text: %{body: message}
    }

    headers = [
      {"Authorization", "Bearer #{config.access_token}"},
      {"Content-Type", "application/json"}
    ]

    case HTTPoison.post(url, Jason.encode!(body), headers) do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        Logger.info("Mensagem WhatsApp enviada com sucesso",
          phone_number: phone_number,
          response: body
        )

        {:ok, Jason.decode!(body)}

      {:ok, %HTTPoison.Response{status_code: code, body: body}} ->
        Logger.error("Erro ao enviar mensagem WhatsApp",
          phone_number: phone_number,
          status_code: code,
          response: body
        )

        {:error, body}

      {:error, %HTTPoison.Error{reason: reason}} ->
        Logger.error("Falha na requisição WhatsApp",
          phone_number: phone_number,
          reason: inspect(reason)
        )

        {:error, :request_failed}
    end
  end

  defp get_config do
    case {System.get_env("WHATSAPP_PHONE_NUMBER_ID"), System.get_env("WHATSAPP_ACCESS_TOKEN")} do
      {nil, _} ->
        Logger.error("Variável de ambiente WHATSAPP_PHONE_NUMBER_ID não configurada")
        {:error, :missing_phone_number_id}

      {_, nil} ->
        Logger.error("Variável de ambiente WHATSAPP_ACCESS_TOKEN não configurada")
        {:error, :missing_access_token}

      {phone_number_id, access_token} ->
        {:ok, %{phone_number_id: phone_number_id, access_token: access_token}}
    end
  end

  defp validate_phone_number(phone_number) do
    if Regex.match?(~r/^\d{10,15}$/, phone_number) do
      {:ok, phone_number}
    else
      Logger.warning("Número de telefone inválido", phone_number: phone_number)
      {:error, :invalid_phone_number}
    end
  end

  defp validate_message(message) when byte_size(message) > 0 and byte_size(message) <= 4096 do
    {:ok, message}
  end

  defp validate_message(_message) do
    Logger.warning("Mensagem inválida: deve ter entre 1 e 4096 caracteres")
    {:error, :invalid_message}
  end
end
