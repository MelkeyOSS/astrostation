import { useSpotifyMusic } from "@Store";
import { useState } from "react";
import { AiOutlineReload } from "react-icons/ai";
import { IoCloseSharp } from "react-icons/io5";
import { WithTooltip } from "../../Tooltip";

export const Spotify = () => {
  const { setIsSpotifyToggled } = useSpotifyMusic();

  const [text, setText] = useState("");
  const [playlist, setPlaylist] = useState(
    "https://open.spotify.com/embed/playlist/37i9dQZF1DWWQRwui0ExPn"
  );

  function changePlaylist() {
    if (!text.includes("https://open.spotify.com/playlist/")) {
      alert("Invalid spotify URL");
      return;
    }
    const splitOn = (slicable: string, ...indices: number[]) =>
      [0, ...indices].map((n, i, m) => slicable.slice(n, m[i + 1]));
    const stitchUrl = splitOn(text, 24)[0] + "/embed" + splitOn(text, 24)[1];
    setPlaylist(stitchUrl);
    setText("");
  }

  function handleKeyDown(e) {
    if (e.key === "Enter") {
      changePlaylist();
    }
  }

  return (
    <div className="py-2 mb-2 w-72 sm:w-96 max-w-sm text-gray-800 shadow-md rounded-lg dark:text-gray-300 bg-white/[.96] dark:bg-gray-800/[.96] dark:border-gray-700 justify-between">
      <WithTooltip text="Make sure to refresh after logging in">
        <div className="flex justify-between items-center p-1 handle cursor-move">
          <p>Spotify</p>
          <IoCloseSharp
            className="text-red-500 cursor-pointer hover:bg-red-200"
            onClick={() => setIsSpotifyToggled(false)}
          />
        </div>
      </WithTooltip>

      <div className="cancelDrag justify-center">
        <iframe
          src={`${playlist}?utm_source=generator&theme=0`}
          height="380"
          width="100%"
          frameBorder="0"
          allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
        ></iframe>
      </div>
      <div className="flex items-center space-x-1 p-1">
        <input
          className="cancelDrag w-full p-1 border border-gray-300 dark:bg-gray-700/[.96] dark:border-gray-500"
          type="text"
          value={text}
          placeholder="Ctrl-V Spotify URL here"
          onChange={(e) => {
            setText(e.target.value);
          }}
          onKeyDown={handleKeyDown}
        />
        <AiOutlineReload
          className="w-5 cursor-pointer hover:text-slate-500"
          onClick={changePlaylist}
        />
      </div>
    </div>
  );
};
